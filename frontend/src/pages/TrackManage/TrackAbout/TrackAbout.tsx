import { useState } from 'react';
import {
	Alert,
	Avatar,
	Box,
	Button,
	Group,
	Stack,
	Text,
	Textarea,
	TextInput,
	Title,
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { UpdateTrackPayload } from '@/api/services/workouts';
import { Track } from '@/types/domain';
import classes from './TrackAbout.module.css';

export interface TrackAboutProps {
	track: Track;

	/** Резолвится, когда трек сохранён; отказ оставляет форму открытой. */
	onSave: (payload: UpdateTrackPayload) => Promise<unknown>;
	isSaving?: boolean;
	error?: Error | null;
}

/**
 * TrackAbout — тексты трека с правкой на месте.
 *
 * Правится там же, где читается: владелец видит трек ровно таким, каким его
 * увидят участники, и не уходит за этим на отдельный экран настроек.
 */
export function TrackAbout({ track, onSave, isSaving = false, error }: TrackAboutProps) {
	const [editing, setEditing] = useState(false);

	const form = useForm<UpdateTrackPayload>({
		mode: 'uncontrolled',
		initialValues: { Name: track.Name, Description: track.Description },
		validate: {
			Name: isNotEmpty('Название не может быть пустым'),
		},
		enhanceGetInputProps: () => ({
			disabled: isSaving,
		}),
	});

	function startEditing() {
		form.initialize({ Name: track.Name, Description: track.Description });
		setEditing(true);
	}

	function cancelEditing() {
		setEditing(false);
	}

	async function handleSubmit(values: UpdateTrackPayload) {
		try {
			await onSave(values);
			setEditing(false);
		} catch {
			// Форма остаётся открытой: набранный текст не теряется, а почему
			// не сохранилось — показано алертом над полями
		}
	}

	return (
		<Group align="flex-start" justify="space-between" gap="md">
			<Box className={classes.about}>
				<Text ff="heading" fz="sm" fw={500} lts="0.12em" tt="uppercase" c="copper.6" mb="xs">
					Управление треком
				</Text>

				{editing ? (
					<form onSubmit={form.onSubmit(handleSubmit)}>
						<Stack gap="md">
							{error && (
								<Alert color="red" variant="light">
									Не удалось сохранить. Попробуйте ещё раз.
								</Alert>
							)}

							<TextInput
								label="Название трека"
								key={form.key('Name')}
								{...form.getInputProps('Name')}
							/>
							<Textarea
								label="Описание"
								autosize
								minRows={4}
								key={form.key('Description')}
								{...form.getInputProps('Description')}
							/>

							<Group gap="md">
								<Button type="submit" loading={isSaving}>
									Сохранить
								</Button>
								<Button variant="default" onClick={cancelEditing} disabled={isSaving}>
									Отмена
								</Button>
								<Text span fz="sm" c="gray.7">
									Текст виден всем участникам трека
								</Text>
							</Group>
						</Stack>
					</form>
				) : (
					<>
						<Title order={1} size="h2" mb="xs">
							{track.Name}
						</Title>
						<Text fz="md" lh="md" c="gray.8" textWrap="pretty" className={classes.description}>
							{track.Description}
						</Text>
					</>
				)}

				{track.Author?.Name && (
					<Group gap="xs" mt="lg">
						<Avatar size={26} radius="xl" name={track.Author.Name} color="copper" />
						<Text span fz="sm" c="gray.8">
							{track.Author.Name}
						</Text>
					</Group>
				)}
			</Box>

			{!editing && (
				<Button variant="default" flex="none" onClick={startEditing}>
					Изменить текст
				</Button>
			)}
		</Group>
	);
}
