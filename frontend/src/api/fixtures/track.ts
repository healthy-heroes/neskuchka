import { TrackData } from '@/api/services/workouts';

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
