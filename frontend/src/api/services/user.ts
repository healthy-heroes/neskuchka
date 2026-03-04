import { UseQueryOptions } from '@tanstack/react-query';
import axios from 'axios';
import { User } from '@/types/domain';
import Service from './service';

export const UserKeys = {
	me: ['user', 'me'] as const,
};

type UserResponse = {
	data: User;
};

export class UserService extends Service {
	/**
	 * Get the current user if authenticated, null if not authenticated (401)
	 */
	getUserQuery(): UseQueryOptions<UserResponse | null> {
		return {
			queryKey: UserKeys.me,
			queryFn: async () => {
				try {
					return await this.api.get<UserResponse>(`user/me`);
				} catch (error) {
					if (axios.isAxiosError(error) && error.response?.status === 401) {
						return null;
					}
					throw error;
				}
			},
			staleTime: 5 * 60 * 1000, // 5 minutes
		};
	}
}
