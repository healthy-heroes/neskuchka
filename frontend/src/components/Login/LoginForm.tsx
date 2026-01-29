import { IconAlertCircle, IconMail } from '@tabler/icons-react';
import { Alert, Button, Container, Paper, Stack, Text, TextInput, Title } from '@mantine/core';
import { isEmail, isNotEmpty, useForm } from '@mantine/form';

type LoginFormData = {
	email: string;
};

export type LoginFormProps = {
	onSubmit: (email: string) => void;
	isSubmitting?: boolean;
	error: Error | null;
};

const emailValidators = [
	isNotEmpty('Электронная почта обязательна'),
	isEmail('Некорректный email'),
];

export function LoginForm({ onSubmit, isSubmitting, error }: LoginFormProps) {
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
		<Container size={460} my={40}>
			<Stack gap="lg" align="center" mb="xl">
				<Title order={1} ta="center" size="h1">
					Войти в Нескучку
				</Title>
				<Text size="sm" c="dimmed" ta="center">
					Введите свой email, и мы отправим вам ссылку для входа
				</Text>
			</Stack>

			{error && (
				<Alert variant="light" color="red" icon={<IconAlertCircle size={20} />} mb="lg">
					{error.message}
				</Alert>
			)}

			<Paper withBorder shadow="md" p={40} radius="md">
				<form onSubmit={form.onSubmit(handleSubmit)}>
					<Stack gap="lg">
						<TextInput
							size="md"
							leftSection={<IconMail size={18} />}
							label="Электронная почта"
							placeholder="your@email.com"
							{...form.getInputProps('email')}
						/>
						<Button type="submit" size="md" fullWidth loading={isSubmitting}>
							Отправить ссылку для входа
						</Button>
					</Stack>
				</form>
			</Paper>
		</Container>
	);
}
