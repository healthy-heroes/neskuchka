import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import axios from 'axios';
import { User, UserSettings } from '@/types/domain';
import Service from './service';

export const UserKeys = {
	me: ['user', 'me'] as const,
	settings: ['user', 'settings'] as const,
};

type UserResponse = {
	data: User;
};

type SettingsResponse = {
	data: UserSettings;
};

export class UserService extends Service {
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
			staleTime: 5 * 60 * 1000,
		};
	}

	getSettingsQuery(): UseQueryOptions<SettingsResponse, Error, UserSettings> {
		return {
			queryKey: UserKeys.settings,
			queryFn: () => this.api.get<SettingsResponse>(`user/me/settings`),
			select: (response) => response.data,
		};
	}

	updateSettingsMutation(): UseMutationOptions<UserSettings, Error, { Name: string }> {
		return {
			mutationFn: async (payload) => {
				const response = await this.api.put<SettingsResponse>(`user/me/settings`, payload);
				return response.data;
			},
		};
	}

	uploadAvatarMutation(): UseMutationOptions<unknown, Error, File> {
		return {
			mutationFn: async (file) => {
				const formData = new FormData();
				formData.append('avatar', file);
				return this.api.postForm(`user/me/avatar`, formData);
			},
		};
	}

	deleteAvatarMutation(): UseMutationOptions<unknown, Error, void> {
		return {
			mutationFn: async () => {
				return this.api.delete(`user/me/avatar`);
			},
		};
	}
}
