import { useEffect } from 'react';
import { IconAlertCircle, IconCheck } from '@tabler/icons-react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Container, Loader, Stack, Text } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { RouteLink } from '@/components/RouteLink/RouteLink';

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
		return (
			<Container size={600} my={30}>
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
			</Container>
		);
	}

	if (!token) {
		return (
			<Container size={600} my={30}>
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
			</Container>
		);
	}

	if (confirmMutation.error) {
		return (
			<Container size={600} my={30}>
				<Alert
					icon={<IconAlertCircle size={24} />}
					title="Ошибка входа"
					color="red"
					variant="light"
					p="xl"
				>
					<Stack gap="md" mt="xs">
						<Text size="md">{confirmMutation.error.message}</Text>
						<RouteLink to="/login" size="sm">
							← Вернуться к форме входа
						</RouteLink>
					</Stack>
				</Alert>
			</Container>
		);
	}

	return (
		<Container size={600} my={30}>
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
		</Container>
	);
}
