import dayjs from 'dayjs';
import {
	IconSquareRoundedMinus2,
	IconSquareRoundedPlus,
	IconSquareRoundedPlus2,
	IconTrash,
} from '@tabler/icons-react';
import {
	ActionIcon,
	Alert,
	Button,
	Card,
	Divider,
	Fieldset,
	Group,
	Textarea,
	TextInput,
	Title,
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { formRootRule, isNotEmpty, useForm, UseFormReturnType } from '@mantine/form';
import { randomId } from '@mantine/hooks';
import { Workout, WorkoutExercise, WorkoutSection } from '@/types/domain';

type WorkoutFormProps = {
	data?: Workout;
	onSubmit: (values: Workout) => void;
	onCancel: () => void;

	isSubmitting?: boolean;
	error: Error | null;
};

type FormReturnType = UseFormReturnType<Workout>;

/**
 * WorkoutForm - component for creating and editing workout
 */
export function WorkoutForm({
	data,
	onSubmit,
	onCancel,
	isSubmitting = false,
	error,
}: WorkoutFormProps) {
	const form = useForm<Workout>({
		mode: 'uncontrolled',
		initialValues: data ?? makeInitialValues(),
		enhanceGetInputProps: () => {
			if (isSubmitting) {
				return { disabled: true };
			}

			return {};
		},
		validate: {
			Date: isNotEmpty('Дата тренировки обязательна'),
			Sections: {
				[formRootRule]: isNotEmpty('Хотя бы одна секция обязательна'),
				Protocol: {
					Title: isNotEmpty(),
				},
				Exercises: {
					[formRootRule]: isNotEmpty('Хотя бы одно упражнение обязательно'),
					Description: isNotEmpty(),
				},
			},
		},
	});

	function handleSubmit(values: Workout) {
		onSubmit(values);
	}

	return (
		<div>
			<form onSubmit={form.onSubmit(handleSubmit)}>
				<DatePickerInput
					withAsterisk
					valueFormat="DD MMMM YYYY"
					label="Дата тренировки"
					placeholder="Выберите дату"
					key={form.key('Date')}
					{...form.getInputProps('Date')}
				/>

				{renderSections(form)}

				<Textarea
					mt="md"
					label="Примечания"
					placeholder="Заметки к тренировке"
					autosize
					minRows={3}
					key={form.key('Notes')}
					{...form.getInputProps('Notes')}
				/>

				<Divider my="md" />
				{error && (
					<Alert mb="md" color="red">
						{error.message}
					</Alert>
				)}
				<Group justify="space-between">
					<Button color="green.7" type="submit" loading={isSubmitting}>
						Сохранить
					</Button>
					<Button color="red.7" onClick={onCancel} disabled={isSubmitting}>
						Отменить
					</Button>
				</Group>
			</form>
		</div>
	);
}

// Additional render helpers
function renderSections(form: FormReturnType): React.ReactNode {
	return (
		<>
			<Title order={2} mt="md">
				Секции
			</Title>
			{renderSectionsFields(form)}
			<Button
				leftSection={<IconSquareRoundedPlus2 size={16} />}
				variant="outline"
				onClick={() => addSection(form)}
				{...form.getInputProps('Buttons.addSection')}
			>
				Добавить секцию
			</Button>
		</>
	);
}

function renderSectionsFields(form: FormReturnType): React.ReactNode {
	return form.getValues().Sections.map((_, index) => {
		const pathFor = (property: string) => `Sections.${index}.${property}`;
		const keyBy = (property: string) => form.key(pathFor(property));

		return (
			<Card key={randomId()} withBorder mb="md">
				<TextInput
					label="Название секции"
					placeholder="Разминка, Комплекс, ... "
					withAsterisk
					style={{ flex: 1 }}
					key={keyBy('Title')}
					{...form.getInputProps(pathFor('Title'))}
				/>

				<Fieldset mt="xs" legend="Протокол">
					<TextInput
						label="Название"
						placeholder="EMOM 20 минут"
						withAsterisk
						style={{ flex: 1 }}
						key={keyBy('Protocol.Title')}
						{...form.getInputProps(pathFor('Protocol.Title'))}
					/>
					<TextInput
						label="Пояснения"
						placeholder="Например: В начале каждой минуты выполнять упражнения"
						style={{ flex: 1 }}
						key={keyBy('Protocol.Description')}
						{...form.getInputProps(pathFor('Protocol.Description'))}
					/>
				</Fieldset>

				{renderExercises(form, index)}

				<Group justify="flex-end">
					<Button
						mt="md"
						leftSection={<IconSquareRoundedMinus2 size={16} />}
						variant="outline"
						color="red"
						onClick={() => removeSection(form, index)}
						{...form.getInputProps('Buttons.removeSection')}
					>
						Удалить секцию
					</Button>
				</Group>
			</Card>
		);
	});
}

function renderExercises(form: FormReturnType, sectionIndex: number): React.ReactNode {
	return (
		<Fieldset mt="md" legend="Упражнения">
			{renderExercisesFields(form, sectionIndex)}

			<Button
				leftSection={<IconSquareRoundedPlus size={16} />}
				variant="outline"
				onClick={() => addExercise(form, sectionIndex)}
				{...form.getInputProps('Buttons.addExercise')}
			>
				Добавить упражнение
			</Button>
		</Fieldset>
	);
}

function renderExercisesFields(form: FormReturnType, sectionIndex: number): React.ReactNode {
	return form.getValues().Sections[sectionIndex].Exercises.map((_, index) => {
		const pathFor = (property: string) => `Sections.${sectionIndex}.Exercises.${index}.${property}`;
		const keyBy = (property: string) => form.key(pathFor(property));

		return (
			<Group mb="xs" key={randomId()}>
				<TextInput
					placeholder="Описание упражнения"
					withAsterisk
					style={{ flex: 1 }}
					key={keyBy('Description')}
					{...form.getInputProps(pathFor('Description'))}
				/>
				<ActionIcon
					color="red"
					onClick={() => removeExercise(form, sectionIndex, index)}
					{...form.getInputProps('Buttons.removeExercise')}
				>
					<IconTrash size={16} />
				</ActionIcon>
			</Group>
		);
	});
}

// Helpers for adding and removing items
function addSection(form: FormReturnType) {
	form.insertListItem('Sections', makeSection());
}

function removeSection(form: FormReturnType, sectionIndex: number) {
	if (form.getValues().Sections.length === 1) {
		form.setFieldValue('Sections', [makeSection()]);
	} else {
		form.removeListItem('Sections', sectionIndex);
	}
}

function addExercise(form: FormReturnType, sectionIndex: number) {
	form.insertListItem(`Sections.${sectionIndex}.Exercises`, makeExercise());
}

function removeExercise(form: FormReturnType, sectionIndex: number, exerciseIndex: number) {
	if (form.getValues().Sections[sectionIndex].Exercises.length === 1) {
		form.setFieldValue(`Sections.${sectionIndex}.Exercises`, [makeExercise()]);
	} else {
		form.removeListItem(`Sections.${sectionIndex}.Exercises`, exerciseIndex);
	}
}

// Helpers for creating initial values
function makeInitialValues(): Workout {
	return {
		ID: randomId('new'),
		Date: dayjs().format('YYYY-MM-DD'),
		Sections: [makeSection('Разминка'), makeSection('Комплекс')],
	};
}

function makeSection(title: string = 'Комплекс'): WorkoutSection {
	return {
		Title: title,
		Protocol: {
			Title: '',
			Description: '',
		},
		Exercises: [makeExercise()],
	};
}

function makeExercise(): WorkoutExercise {
	return {
		ExerciseSlug: randomId('new'),
		Description: '',
	};
}
