import { useState } from 'react';
import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Box, Button, Container, Group, Modal, Stack, Text, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { MANAGE_PAGE_SIZE, WorkoutsKeys } from '@/api/services/workouts';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { RouteLink } from '@/components/RouteLink/RouteLink';
import { Workout } from '@/types/domain';
import { formatIsoDateShort } from '@/utils/dates';
import { TrackAbout } from './TrackAbout/TrackAbout';
import { WorkoutRows } from './WorkoutRows/WorkoutRows';
import classes from './TrackManage.module.css';

/**
 * TrackManagePage — трек глазами владельца: тексты, все тренировки и действия
 * над ними на одном экране.
 *
 * Права уже проверены гардом TrackOwnerOnly, страница их не перепроверяет.
 */
export function TrackManagePage() {
	const queryClient = useQueryClient();
	const { workouts } = useApi();

	const [confirming, setConfirming] = useState<Workout | null>(null);

	const trackQuery = useQuery(workouts.getMainTrackQuery());
	const listQuery = useInfiniteQuery(workouts.getAllTrackWorkoutsQuery());

	// Тренировка есть в двух списках сразу — в админском и в публичном, — а трек
	// ещё и в шапке страницы, так что после любой правки сбрасывается вся ветка
	const invalidateTrack = () => queryClient.invalidateQueries({ queryKey: WorkoutsKeys.track() });

	const updateTrack = useMutation({
		...workouts.updateTrackMutation(),
		onSuccess: invalidateTrack,
	});

	const deleteWorkout = useMutation({
		...workouts.deleteWorkoutMutation(),
		onSuccess: invalidateTrack,
	});

	if (trackQuery.isPending || listQuery.isPending) {
		return <PageSkeleton />;
	}

	if (!trackQuery.data || !listQuery.data) {
		return null;
	}

	const pages = listQuery.data.pages;
	const rows = pages.flatMap((page) => page.data.Workouts);
	const { Total, Planned } = pages[0].data;

	function handleConfirmDelete() {
		if (!confirming) {
			return;
		}

		deleteWorkout.mutate(confirming.ID, { onSuccess: () => setConfirming(null) });
	}

	return (
		<Container className={classes.page}>
			<Box className={classes.about}>
				<TrackAbout
					track={trackQuery.data.Track}
					onSave={updateTrack.mutateAsync}
					isSaving={updateTrack.isPending}
					error={updateTrack.error}
				/>
			</Box>

			<Group align="flex-end" justify="space-between" gap="xl" py="xl">
				<div>
					<Title order={3}>Тренировки</Title>
					<Text fz="sm" c="gray.7" mt={6}>
						{countLabel(Total, Planned)}
					</Text>
				</div>

				<Stack align="flex-end" gap={6}>
					<Button component={RouteLink} to="/workouts/new" underline="never">
						Новая тренировка
					</Button>
					<Text fz="sm" c="gray.7">
						Дата — сегодня или позже
					</Text>
				</Stack>
			</Group>

			{rows.length === 0 ? (
				<EmptyTrack />
			) : (
				<>
					<Group gap="xl" mb="sm" px="lg" py="sm" bg="gray.1" className={classes.rules}>
						<Text fz="sm" lh="sm" c="gray.7">
							Тренировка появляется у участников в свой день
						</Text>
						<Text fz="sm" lh="sm" c="gray.7">
							Изменить и удалить можно, пока тренировке не больше суток
						</Text>
					</Group>

					<WorkoutRows workouts={rows} onDelete={setConfirming} />

					{listQuery.hasNextPage && (
						<Stack align="center" gap="xs" pt="lg">
							<Button
								variant="default"
								loading={listQuery.isFetchingNextPage}
								onClick={() => listQuery.fetchNextPage()}
							>
								{moreLabel(Total - rows.length)}
							</Button>
							<Text fz="sm" c="gray.7">
								Показаны {rows.length} из {Total}
							</Text>
						</Stack>
					)}
				</>
			)}

			<Modal
				opened={confirming !== null}
				onClose={() => setConfirming(null)}
				title="Удалить тренировку?"
				centered
			>
				<Text fz="md" lh="md" c="gray.8">
					{confirming &&
						`Тренировка ${formatIsoDateShort(confirming.Date)} исчезнет у всех участников трека. Отменить это нельзя.`}
				</Text>

				<Group justify="flex-end" gap="md" mt="xl">
					<Button
						variant="default"
						disabled={deleteWorkout.isPending}
						onClick={() => setConfirming(null)}
					>
						Отмена
					</Button>
					<Button color="copper.7" loading={deleteWorkout.isPending} onClick={handleConfirmDelete}>
						Удалить
					</Button>
				</Group>
			</Modal>
		</Container>
	);
}

function EmptyTrack() {
	return (
		<Stack align="center" gap="md" px="xl" py={64} bg="white" className={classes.empty}>
			<Title order={3}>В треке пока нет тренировок</Title>
			<Text fz="md" lh="md" c="gray.7" ta="center" maw={420}>
				Создайте первую — участники увидят её в тот день, на который она поставлена.
			</Text>
			<Button component={RouteLink} to="/workouts/new" underline="never" mt="xs">
				Новая тренировка
			</Button>
		</Stack>
	);
}

function countLabel(total: number, planned: number): string {
	if (total === 0) {
		return 'Ни одной тренировки';
	}

	return `Всего в треке ${total} ${plural(total, ['тренировка', 'тренировки', 'тренировок'])} · ${planned} ${plural(planned, ['запланирована', 'запланированы', 'запланировано'])}`;
}

function moreLabel(rest: number): string {
	const next = Math.min(rest, MANAGE_PAGE_SIZE);

	return `Показать ещё ${next} ${plural(next, ['тренировку', 'тренировки', 'тренировок'])}`;
}

/** Русский счётный род: одна тренировка, две тренировки, пять тренировок. */
function plural(count: number, forms: [string, string, string]): string {
	const ones = count % 10;
	const tens = count % 100;

	if (ones === 1 && tens !== 11) {
		return forms[0];
	}

	if (ones >= 2 && ones <= 4 && (tens < 10 || tens >= 20)) {
		return forms[1];
	}

	return forms[2];
}
