import { useEffect } from 'react';
import { IconAlertCircle, IconCheck } from '@tabler/icons-react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Container, Loader, Stack, Text } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { AuthKeys } from '@/api/services/auth';
import { RouteLink } from '@/components/RouteLink/RouteLink';

export type LoginConfirmProps = {
	token?: string;
};

export function LoginConfirm({ token }: LoginConfirmProps) {
	const { auth } = useApi();
	const queryClient = useQueryClient();
	const confirmQuery = useQuery({
		...auth.confirmLoginQuery(token || ''),
		enabled: !!token,
	});

	// Invalidate user query on success
	useEffect(() => {
		if (confirmQuery.isSuccess) {
			queryClient.invalidateQueries({ queryKey: AuthKeys.user });
		}
	}, [confirmQuery.isSuccess, queryClient]);

	// Helper for wrapping content in a container
	function render(children: React.ReactNode) {
		return (
			<Container size={600} my={30}>
				{children}
			</Container>
		);
	}

	// No token state
	if (!token) {
		return render(<NoTokenState />);
	}

	// Token confirmed by backend
	if (confirmQuery.isSuccess) {
		return render(<SuccessState />);
	}

	// Request completed with error
	if (confirmQuery.isError) {
		return render(<ErrorState message={confirmQuery.error.message} />);
	}

	return render(<LoadingState />);
}

function SuccessState() {
	return (
		<Alert
			icon={<IconCheck size={24} />}
			title="Вход выполнен успешно!"
			color="green"
			variant="light"
			p="xl"
		>
			<Stack gap="md" mt="xs">
				<Text size="md">Вы успешно вошли в систему. Готовы начать тренироваться?</Text>
				<RouteLink to="/workouts" size="sm">
					Начать →
				</RouteLink>
			</Stack>
		</Alert>
	);
}

function NoTokenState() {
	return (
		<Alert
			icon={<IconAlertCircle size={24} />}
			title="Отсутствует токен"
			color="red"
			variant="light"
			p="xl"
		>
			<Stack gap="md" mt="xs">
				<Text size="md">
					Ссылка для входа недействительна или устарела. Пожалуйста, запросите новую ссылку.
				</Text>
				<RouteLink to="/login" size="sm">
					← Вернуться к форме входа
				</RouteLink>
			</Stack>
		</Alert>
	);
}

function ErrorState({ message }: { message: string }) {
	return (
		<Alert
			icon={<IconAlertCircle size={24} />}
			title="Ошибка входа"
			color="red"
			variant="light"
			p="xl"
		>
			<Stack gap="md" mt="xs">
				<Text size="md">{message}</Text>
				<RouteLink to="/login" size="sm">
					← Вернуться к форме входа
				</RouteLink>
			</Stack>
		</Alert>
	);
}

function LoadingState() {
	return (
		<Alert
			icon={<Loader size={24} />}
			title="Подтверждение входа..."
			color="blue"
			variant="light"
			p="xl"
		>
			<Text size="md" mt="xs">
				Пожалуйста, подождите, мы проверяем вашу ссылку...
			</Text>
		</Alert>
	);
}
