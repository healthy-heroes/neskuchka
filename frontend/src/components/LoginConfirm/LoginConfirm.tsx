import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useApi } from '@/api/hooks';

export type LoginConfirmProps = {
	token?: string;
};

export function LoginConfirm({ token }: LoginConfirmProps) {
	const { auth } = useApi();
	const confirmMutation = useMutation(auth.confirmLoginMutation());

	useEffect(() => {
		if (token && !confirmMutation.isSuccess && !confirmMutation.isPending) {
			confirmMutation.mutate(token);
		}
	}, [token]);

	if (confirmMutation.isSuccess) {
		return <div>Login confirmed</div>;
	}

	if (!token) {
		return <div>No token</div>;
	}

	if (confirmMutation.error) {
		return <div>Error: {confirmMutation.error.message}</div>;
	}

	return <div>Confirming...</div>;
}
