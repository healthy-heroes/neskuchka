import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Container, Group, Paper, TextInput, Title } from '@mantine/core';
import { isEmail, isNotEmpty, useForm } from '@mantine/form';
import { useApi } from '@/api/hooks';

export function Login() {
	const { auth } = useApi();

	const loginMutation = useMutation(auth.loginMutation());

	if (loginMutation.isSuccess) {
		return <div>Email sent</div>;
	}

	return (
		<LoginForm
			isSubmitting={loginMutation.isPending}
			onSubmit={loginMutation.mutate}
			error={loginMutation.error}
		/>
	);
}

export function LoginConfirm({ token }: { token?: string }) {
	const { auth } = useApi();
	const [isSubmitting, setIsSubmitting] = useState(false);

	const confirmMutation = useMutation(auth.confirmLoginMutation());

	if (confirmMutation.isSuccess) {
		return <div>Login confirmed</div>;
	}

	if (!token) {
		return <div>No token</div>;
	}

	if (confirmMutation.error) {
		return <div>Error: {confirmMutation.error.message}</div>;
	}

	if (!isSubmitting) {
		setIsSubmitting(true);
		setTimeout(() => {
			confirmMutation.mutate(token);
		}, 1);
	}

	return <div>Confirming...</div>;
}

type LoginFormData = {
	email: string;
};

type LoginFormProps = {
	onSubmit: (email: string) => void;
	isSubmitting?: boolean;
	error: Error | null;
};

const emailValidators = [
	isNotEmpty('Электронная почта обязательна'),
	isEmail('Некорректный email'),
];

function LoginForm({ onSubmit, isSubmitting, error }: LoginFormProps) {
	const form = useForm<LoginFormData>({
		mode: 'uncontrolled',
		initialValues: {
			email: '',
		},
		enhanceGetInputProps: () => {
			if (isSubmitting) {
				return { disabled: true };
			}

			return {};
		},
		validate: {
			email: (value) => {
				const errors = emailValidators.map((v) => v(value)).filter(Boolean);
				return errors.length > 0 ? errors[0] : null;
			},
		},
	});

	function handleSubmit(values: LoginFormData) {
		onSubmit(values.email);
	}

	return (
		<Container size={460} my={30}>
			<Title ta="center">Войти в Нескучку</Title>
			{error && (
				<Alert mt="xl" color="red">
					{error.message}
				</Alert>
			)}
			<Paper withBorder shadow="md" p={30} radius="md" mt="xl">
				<form onSubmit={form.onSubmit(handleSubmit)}>
					<TextInput
						withAsterisk
						label="Электронная почта"
						placeholder="some@neskuchka.ru"
						{...form.getInputProps('email')}
					/>
					<Group justify="space-between" mt="lg">
						<Button type="submit" loading={isSubmitting}>
							Войти
						</Button>
					</Group>
				</form>
			</Paper>
		</Container>
	);
}
