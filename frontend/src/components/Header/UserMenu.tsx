import { IconAdjustments, IconLogout, IconSettings } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Avatar, Button, Menu, Skeleton } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { useAuth } from '@/auth/hooks';
import { RouteLink } from '../RouteLink/RouteLink';

export function UserMenu() {
	const { user, isAuthenticated, isPending, logout } = useAuth();

	// Шапка живёт на всех страницах, так что на /workouts трек берётся из кэша,
	// а вне его — это один лишний запрос. Он публичный и гостю отвечает 200,
	// так что дешевле запросить, чем городить условие
	const { workouts } = useApi();
	const { data: track } = useQuery(workouts.getMainTrackQuery());

	if (isPending) {
		return <Skeleton height={36} width={36} circle />;
	}

	if (isAuthenticated && user) {
		return (
			<Menu shadow="md" width={200}>
				<Menu.Target>
					<Avatar
						src={user.Avatar}
						name={user.Name}
						color="copper"
						size="md"
						radius="xl"
						style={{ cursor: 'pointer' }}
					/>
				</Menu.Target>

				<Menu.Dropdown>
					<Menu.Label>{user.Name}</Menu.Label>
					{track?.IsOwner && (
						<Menu.Item
							component={RouteLink}
							to="/workouts/manage"
							leftSection={<IconAdjustments size={16} />}
						>
							Управление треком
						</Menu.Item>
					)}
					<Menu.Item component={RouteLink} to="/settings" leftSection={<IconSettings size={16} />}>
						Настройки
					</Menu.Item>
					<Menu.Divider />
					<Menu.Item color="red" leftSection={<IconLogout size={16} />} onClick={() => logout()}>
						Выйти
					</Menu.Item>
				</Menu.Dropdown>
			</Menu>
		);
	}

	return (
		<Button component={RouteLink} to="/login">
			Войти
		</Button>
	);
}
