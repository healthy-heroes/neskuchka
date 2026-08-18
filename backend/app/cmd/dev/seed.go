package devcmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/cmd"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/datastorage"
	"github.com/healthy-heroes/neskuchka/backend/app/storage/db"
)

type SeedCommand struct {
	Store cmd.StoreOptions `group:"store" namespace:"store" env-namespace:"STORE"`

	cmd.CommonOptions
}

type SeedRunner struct {
	engine      *db.Engine
	dataStorage *datastorage.Storage
}

func (cmd *SeedCommand) Execute(args []string) error {
	log.Info().Msg("running seed...")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn().Msg("got interrupt signal")
		cancel()
	}()

	runner, err := cmd.createRunner()
	if err != nil {
		return fmt.Errorf("failed create seed runner: %w", err)
	}

	if err := runner.Run(ctx); err != nil {
		return fmt.Errorf("failed execute runner: %w", err)
	}

	log.Info().Msg("database seeded successfully")
	return nil
}

func (cmd *SeedCommand) createRunner() (*SeedRunner, error) {
	log.Info().Msg("creating store")
	log.Info().Msgf("database url: %s", cmd.Store.DB)

	engine, err := db.NewEngine(cmd.Store.DB, db.Opts{Logger: log.Logger})
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	return &SeedRunner{
		engine:      engine,
		dataStorage: datastorage.New(engine, log.Logger),
	}, nil
}

func (r *SeedRunner) Run(ctx context.Context) error {
	log.Info().Msg("start seed runner")

	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Info().Msg("runner shutdown...")
		if err := r.engine.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database")
		}
	}()

	admin := domain.User{
		ID:    domain.NewUserID(),
		Name:  "Илья Карягин",
		Email: "admin@example.com",
	}

	// check existing admin
	_, err := r.dataStorage.GetUserByEmail(ctx, admin.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("failed to check existing data: %w", err)
	}
	if err == nil {
		log.Info().Msg("data already exists, skipping")
		return nil
	}

	// Create user
	_, err = r.dataStorage.CreateUser(ctx, admin)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create track
	mainTrack := domain.Track{
		ID:          domain.NewTrackID(),
		Name:        "Нескучный спорт",
		Slug:        domain.TrackSlug("main"),
		Description: "Тренируйтесь с нами — где бы вы ни находились!\nИдеальная программа, чтобы поддерживать форму дома, без специального оборудования.",
		OwnerID:     admin.ID,
	}
	_, err = r.dataStorage.CreateTrack(ctx, mainTrack)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	for _, w := range seedWorkouts(mainTrack.ID, time.Now()) {
		if _, err := r.dataStorage.CreateWorkout(ctx, w); err != nil {
			return fmt.Errorf("failed to create workout: %w", err)
		}
	}

	return nil
}

// ex builds an exercise. Several prescriptions mean several sets, each on its own line.
// A warm-up exercise often has no prescription at all; keep the slice non-nil so the
// stored JSON holds an empty array rather than null.
func ex(name string, prescription ...string) domain.WorkoutExercise {
	if prescription == nil {
		prescription = []string{}
	}

	return domain.WorkoutExercise{Name: name, Prescription: prescription}
}

// section builds a workout section. Note goes under the section heading and may be empty.
// Exercises are required: a section without them is not a state the product allows,
// and the editor refuses to save one.
func section(title, protocol, note string, exercises ...domain.WorkoutExercise) domain.WorkoutSection {
	if len(exercises) == 0 {
		panic("seed: section " + title + " has no exercises")
	}

	return domain.WorkoutSection{
		Title: title,
		Protocol: domain.Protocol{
			Type:        domain.ProtocolTypeCustom,
			Title:       protocol,
			Description: note,
		},
		Exercises: exercises,
	}
}

// seedWorkouts builds three weeks of workouts, three per week, counted back from today.
// Dates are relative on purpose: the track page puts today's workout in focus, so the
// seed has to produce one whenever it is run.
//
// Content is taken from the Нескучный спорт telegram channel. There a session is split
// across two posts — [РАЗМИНКА] and [ТРЕНИРОВКА ДНЯ] minutes apart — and both become
// sections of a single workout here.
func seedWorkouts(trackID domain.TrackID, today time.Time) []domain.Workout {
	sections := [][]domain.WorkoutSection{
		{
			section("Разминка", "", "", ex("ПВХ")),
			section("Навык", "", "", ex("Тяжёлая атлетика, толчок (C&J)")),
			section("Комплекс", "", "",
				ex("Швунг жимовой + швунг толчковый (ножницы)", "3х(1+2) @ 70%"),
				ex("Толчок (классика)", "2 @ 70%", "3х2 @ 80%"),
				ex("Присед на груди", "2 @ 70%", "4х4 @ 80%"),
			),
		},
		{
			section("Разминка", "", "",
				ex("Дохлая лягушка"),
				ex("Присед узко"),
				ex("Наклон в сторону ноги крестом"),
			),
			section("Навык", "", "", ex("Метаболическое кондиционирование (MetCon)")),
			section("Комплекс", "21-15-9 на время", "",
				ex("Становая тяга"),
				ex("Фронтальное берпи"),
			),
			section("Дополнительная работа", "", "",
				ex("Жим сидя на наклонной скамье", "3 х 12-15"),
				ex("Тяга в наклоне", "3 х 10+10"),
			),
		},
		{
			section("Разминка", "", "",
				ex("Подъём ног сидя"),
				ex("Качающаяся планка"),
				ex("Присед широко"),
			),
			section("Навык", "Приседание (Back Squat)", "",
				ex("Приседание со штангой на спине", "3-5 х 5 @ 70-80%"),
			),
			section("Комплекс", "3-5 раундов не на время", "",
				ex("Перешагивание коробки", "10+10"),
				ex("Французский жим стоя", "10-15"),
				ex("Фермерская прогулка", "20-30 м"),
			),
		},
		{
			section("Разминка", "", "",
				ex("Ротация грудного отдела сидя"),
				ex("Обратный скорпион"),
				ex("Стол"),
			),
			section("Навык", "", "", ex("Стойка на руках (Hand Stand)")),
			section("Комплекс", "По минутке, 16 минут", "",
				ex("Вис", "30/30 сек"),
				ex("Отжимание в стойке на руках / стойка на руках (20/40)", "5-4-3-2"),
			),
			section("Дополнительная работа", "21-15-9", "Перед комплексом — 1000 м гребли.",
				ex("Отжимание на кольцах / параллетах / полу"),
				ex("Мах гирей"),
			),
		},
		{
			section("Разминка", "", "",
				ex("Краб тач"),
				ex("Гребля", "15/13 кал"),
				ex("Координационная лестница (прыжки)"),
			),
			section("Навык", "", "", ex("Восстановление (Recovery)")),
			section("Комплекс", "3-5 раундов не на время", "",
				ex("Поднос ног", "10-15"),
				ex("Пресс", "10-20"),
				ex("Тяга к лицу (резина)", "10-20"),
				ex("Ягодичный мост с опорой (штанга)", "10"),
			),
			section("Дополнительная работа", "3 раунда не на время", "",
				ex("Присед платформа вверх (босу)", "10"),
				ex("Полуприсед на одной", "10+10"),
				ex("Подъём на носках с возвышения", "10+10"),
			),
		},
		{
			section("Разминка", "3 раунда", "",
				ex("Паук", "10-20"),
				ex("Планка + собака мордой вниз", "10"),
				ex("Снежный ангел", "5"),
			),
			section("Комплекс", "3-5 раундов на время", "Лимит по времени 15 минут. Приседание — с весом, опционально.",
				ex("Берпи", "15"),
				ex("Приседание", "40"),
			),
		},
		{
			section("Разминка", "3 раунда", "",
				ex("Вращение бедра на четвереньках", "5+5+5+5"),
				ex("Отжимание", "10"),
				ex("Прыжок вперёд/назад, вправо/влево", "20+20"),
			),
			section("Комплекс", "5 раундов не на время", "",
				ex("Болгарский выпад", "10+10"),
				ex("Обратное скручивание", "20"),
				ex("Червяк + отжимание", "10-15"),
			),
		},
		{
			section("Разминка", "3 раунда", "",
				ex("Вращение бедра лёжа", "5+5+5+5"),
				ex("Ротация поясницы лёжа", "10+10"),
				ex("Приседание на носках", "10-15"),
			),
			section("Комплекс", "5 раундов не на время", "",
				ex("Марш в планке", "20"),
				ex("Ягодичный марш", "20"),
				ex("Пресс у стены", "20"),
			),
		},
		{
			section("Разминка", "3 раунда", "",
				ex("Обратный скорпион", "10"),
				ex("Раскрытие в обратной планке", "10"),
				ex("Прыжок колени к груди", "3-5"),
			),
			section("Комплекс", "По минутке, 12-20 минут", "В каждом упражнении — максимум повторений.",
				ex("Отжимание"),
				ex("Выпрыгивание из глубокого седа"),
				ex("Складка"),
				ex("Отдых"),
			),
		},
	}

	workouts := make([]domain.Workout, 0, len(sections))
	for i, s := range sections {
		// Three per week, every other day, counted back from today — so the newest
		// workout is always today's: 0, 2, 4, then 7, 9, 11, then 14, 16, 18.
		daysAgo := i/3*7 + i%3*2
		date := today.AddDate(0, 0, -daysAgo)

		workouts = append(workouts, domain.Workout{
			ID:       domain.NewWorkoutID(),
			TrackID:  trackID,
			Date:     time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
			Sections: s,
		})
	}

	return workouts
}
