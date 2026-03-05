import { useEffect } from 'react';
import { IconPhoto, IconTrash, IconUpload } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
	ActionIcon,
	Alert,
	Avatar,
	Button,
	Container,
	FileButton,
	Group,
	Stack,
	Text,
	TextInput,
	Title,
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { useApi } from '@/api/hooks';
import { UserKeys } from '@/api/services/user';
import { Header } from '@/components/Header/Header';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

interface SettingsFormData {
	Name: string;
}

export function SettingsPage() {
	const queryClient = useQueryClient();
	const { user } = useApi();

	const { data: settings, isPending } = useQuery(user.getSettingsQuery());

	const invalidateUserData = () => {
		queryClient.invalidateQueries({ queryKey: UserKeys.me });
		queryClient.invalidateQueries({ queryKey: UserKeys.settings });
	};

	const updateSettings = useMutation({
		...user.updateSettingsMutation(),
		onSuccess: invalidateUserData,
	});
	const uploadAvatar = useMutation({
		...user.uploadAvatarMutation(),
		onSuccess: invalidateUserData,
	});
	const deleteAvatar = useMutation({
		...user.deleteAvatarMutation(),
		onSuccess: invalidateUserData,
	});

	const form = useForm<SettingsFormData>({
		mode: 'uncontrolled',
		initialValues: { Name: '' },
		validate: {
			Name: isNotEmpty('Имя не может быть пустым'),
		},
		enhanceGetInputProps: () => ({
			disabled: updateSettings.isPending,
		}),
	});

	useEffect(() => {
		if (settings) {
			form.initialize({ Name: settings.Name });
		}
	}, [settings]); // eslint-disable-line react-hooks/exhaustive-deps

	if (isPending) {
		return <PageSkeleton />;
	}

	if (!settings) {
		return null;
	}

	function handleSubmit(values: SettingsFormData) {
		updateSettings.mutate({ Name: values.Name });
	}

	function handleAvatarUpload(file: File | null) {
		if (file) {
			uploadAvatar.mutate(file);
		}
	}

	function handleAvatarDelete() {
		deleteAvatar.mutate();
	}

	return (
		<>
			<Header />
			<Container size="xs" py="xl">
				<Title order={2} mb="lg">
					Настройки
				</Title>

				<Stack gap="xl">
					<Stack gap="sm">
						<Text fw={500}>Аватар</Text>
						<Group>
							<Avatar src={settings.Avatar} size={80} radius="xl">
								<IconPhoto size={32} />
							</Avatar>
							<Stack gap="xs">
								<FileButton onChange={handleAvatarUpload} accept="image/png,image/jpeg,image/webp">
									{(props) => (
										<Button
											{...props}
											variant="light"
											size="xs"
											leftSection={<IconUpload size={14} />}
											loading={uploadAvatar.isPending}
										>
											Загрузить
										</Button>
									)}
								</FileButton>
								{settings.Avatar && (
									<ActionIcon
										variant="subtle"
										color="red"
										size="sm"
										onClick={handleAvatarDelete}
										loading={deleteAvatar.isPending}
										title="Удалить аватар"
									>
										<IconTrash size={14} />
									</ActionIcon>
								)}
							</Stack>
						</Group>
						{uploadAvatar.isError && (
							<Alert color="red" variant="light">
								Не удалось загрузить аватар
							</Alert>
						)}
					</Stack>

					<form onSubmit={form.onSubmit(handleSubmit)}>
						<Stack gap="md">
							<TextInput
								label="Email"
								value={settings.Email}
								readOnly
								variant="filled"
							/>
							<TextInput
								label="Имя"
								key={form.key('Name')}
								{...form.getInputProps('Name')}
							/>

							{updateSettings.isError && (
								<Alert color="red" variant="light">
									Не удалось сохранить изменения
								</Alert>
							)}

							<Group>
								<Button type="submit" loading={updateSettings.isPending}>
									Сохранить
								</Button>
							</Group>
						</Stack>
					</form>
				</Stack>
			</Container>
		</>
	);
}
