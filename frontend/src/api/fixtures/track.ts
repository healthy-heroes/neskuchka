import { TrackData, WorkoutsKeys, WorkoutsService } from '@/api/services/workouts';

/** Главный трек — те же данные, что кладёт `mise run //backend:seed`. */
export const mockTrack: TrackData = {
	Track: {
		ID: 'track-1',
		Name: 'Нескучный спорт',
		Description:
			'Тренируйтесь с нами — где бы вы ни находились!\nИдеальная программа, чтобы поддерживать форму дома, без специального оборудования.',
		Author: {
			ID: 'user-1',
			Name: 'Администратор',
		},
	},
	IsOwner: false,
};

/**
 * Мок сервиса треков для стори.
 *
 * Дефолт нужен по той же причине, что у auth и user: шапка живёт почти на
 * каждой странице, а `UserMenu` спрашивает трек, чтобы решить, показывать ли
 * пункт управления. Без дефолта строгий Proxy бросает прямо в рендере, и
 * падает любая стори с `Header` внутри.
 */
export function createWorkoutsServiceMock(
	overrides: Partial<WorkoutsService> = {}
): Partial<WorkoutsService> {
	return {
		getMainTrackQuery: () =>
			({
				queryKey: WorkoutsKeys.track(),
				queryFn: async () => ({ data: mockTrack }),
				select: (response: { data: TrackData }) => response.data,
			}) as ReturnType<WorkoutsService['getMainTrackQuery']>,
		...overrides,
	};
}
