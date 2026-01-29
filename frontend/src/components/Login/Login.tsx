import { useMutation } from '@tanstack/react-query';
import { useApi } from '@/api/hooks';
import { LoginForm } from './LoginForm';
import { LoginSuccess } from './LoginSuccess';

export function Login() {
	const { auth } = useApi();
	const loginMutation = useMutation(auth.loginMutation());

	if (loginMutation.isSuccess) {
		return <LoginSuccess />;
	}

	return (
		<LoginForm
			isSubmitting={loginMutation.isPending}
			onSubmit={loginMutation.mutate}
			error={loginMutation.error}
		/>
	);
}
