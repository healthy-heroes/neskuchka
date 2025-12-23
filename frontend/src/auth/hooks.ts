import {
	UseMutateAsyncFunction,
	useMutation,
	useQuery,
	useQueryClient,
} from '@tanstack/react-query';
import { useApi } from '@/api/hooks';
import { AuthKeys } from '@/api/services/auth';
import { User } from '@/types/domain';

type AuthState = {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	logout: UseMutateAsyncFunction<void, Error, void>;
};

export function useAuth(): AuthState {
	const api = useApi();
	const queryClient = useQueryClient();

	const logout = useMutation({
		...api.auth.logoutMutation(),

		onSuccess: () => {
			queryClient.setQueryData(AuthKeys.user, null);
		},
	});

	const userQuery = useQuery(api.auth.getUserQuery());

	return {
		user: userQuery.data ?? null,
		isAuthenticated: !!userQuery.data,
		isLoading: userQuery.isPending,
		logout: logout.mutateAsync,
	};
}

export function useIsOwner(ownerID?: string): boolean {
	const { user } = useAuth();

	return !!ownerID && !!user && user.ID === ownerID;
}
