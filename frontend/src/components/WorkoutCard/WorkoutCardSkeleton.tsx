import { Card, CardProps, Divider, Grid, Skeleton } from '@mantine/core';

export interface WorkoutCardSkeletonProps {
	cardProps?: CardProps;
}

export function WorkoutCardSkeleton({ cardProps }: WorkoutCardSkeletonProps) {
	return (
		<Card shadow="sm" p="lg" radius="md" withBorder {...cardProps}>
			<Grid>
				<Grid.Col visibleFrom="xs" span={5}>
					<Skeleton height={200} />
				</Grid.Col>
				<Grid.Col span={{ base: 12, xs: 7 }} p="md">
					<Skeleton width="25%" height={30} radius="xl" />

					<Skeleton width="50%" height={8} mt={20} radius="xl" />
					<Skeleton width="45%" height={8} mt={6} radius="xl" />
					<Skeleton width="55%" height={8} mt={6} radius="xl" />

					<Divider my="md" />

					<Skeleton width="50%" height={8} mt={20} radius="xl" />
					<Skeleton width="45%" height={8} mt={6} radius="xl" />
					<Skeleton width="55%" height={8} mt={6} radius="xl" />
				</Grid.Col>
			</Grid>
		</Card>
	);
}
