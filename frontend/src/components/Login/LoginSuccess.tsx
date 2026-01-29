import { IconMailCheck } from '@tabler/icons-react';
import { Alert, Container, List, Stack, Text } from '@mantine/core';
import { RouteLink } from '@/components/RouteLink/RouteLink';

export function LoginSuccess() {
	return (
		<Container size={600} my={30}>
			<Alert
				icon={<IconMailCheck size={24} />}
				title="Письмо отправлено!"
				color="green"
				variant="light"
				p="xl"
			>
				<Stack gap="md" mt="xs">
					<Text size="md">
						Мы отправили вам письмо со ссылкой для входа. Пожалуйста, проверьте свою почту и
						перейдите по ссылке из письма.
					</Text>
					<Text size="sm" c="dimmed" fw={600} mt="xs">
						Что нужно сделать:
					</Text>
					<List type="ordered" size="sm" spacing="xs">
						<List.Item>Откройте письмо от Нескучки в своей почте</List.Item>
						<List.Item>Нажмите на ссылку для подтверждения входа</List.Item>
						<List.Item>Вы будете автоматически перенаправлены обратно</List.Item>
					</List>
					<RouteLink to="/" size="sm">
						← Вернуться на главную страницу
					</RouteLink>
				</Stack>
			</Alert>
		</Container>
	);
}
