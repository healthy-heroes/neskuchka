import { Container, Skeleton, Stack } from '@mantine/core';
import { Header } from '../Header/Header';

export function PageSkeleton() {
	return (
		<>
			<Header />
			<Container size="sm" py="xl">
				<Stack gap="md">
					<Skeleton height={32} width="40%" radius="md" />
					<Skeleton height={16} width="70%" radius="sm" />
					<Skeleton height={16} width="55%" radius="sm" />
					<Skeleton height={200} radius="md" mt="md" />
				</Stack>
			</Container>
		</>
	);
}
