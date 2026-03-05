import { useEffect, useRef } from 'react';
import { IconPencil, IconPhoto, IconTrash } from '@tabler/icons-react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
	Alert,
	Avatar,
	Button,
	Container,
	Group,
	Menu,
	Paper,
	Stack,
	Text,
	TextInput,
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
	const fileInputRef = useRef<HTMLInputElement>(null);

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

	function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (file) {
			uploadAvatar.mutate(file);
		}
		e.target.value = '';
	}

	return (
		<>
			<Header />
			<Container size="xs" py="xl">
				<Stack gap="lg">
					<Paper radius="md" withBorder p="xl">
						<Stack align="center" gap="md">
							<Avatar src={settings.Avatar} size={180} radius={180}>
								{settings.Name.charAt(0).toUpperCase()}
							</Avatar>

							<input
								ref={fileInputRef}
								type="file"
								accept="image/png,image/jpeg,image/webp"
								onChange={handleFileChange}
								hidden
							/>

							<Menu shadow="md" width={200} position="bottom">
								<Menu.Target>
									<Button
										variant="default"
										size="xs"
										leftSection={<IconPencil size={14} />}
										loading={uploadAvatar.isPending || deleteAvatar.isPending}
									>
										Редактировать
									</Button>
								</Menu.Target>
								<Menu.Dropdown>
									<Menu.Item
										leftSection={<IconPhoto size={16} />}
										onClick={() => fileInputRef.current?.click()}
									>
										Загрузить фото...
									</Menu.Item>
									{settings.Avatar && (
										<Menu.Item
											color="red"
											leftSection={<IconTrash size={16} />}
											onClick={() => deleteAvatar.mutate()}
										>
											Удалить фото
										</Menu.Item>
									)}
								</Menu.Dropdown>
							</Menu>

							{uploadAvatar.isError && (
								<Alert color="red" variant="light" w="100%">
									Не удалось загрузить аватар
								</Alert>
							)}
						</Stack>
					</Paper>

					<Paper radius="md" withBorder p="xl">
						<form onSubmit={form.onSubmit(handleSubmit)}>
							<Stack gap="md">
								<Text fw={600} size="lg">
									Личные данные
								</Text>

								<TextInput
									label="Email"
									value={settings.Email}
									readOnly
									variant="filled"
									description="Email нельзя изменить"
								/>
								<TextInput
									label="Имя"
									placeholder="Ваше имя"
									key={form.key('Name')}
									{...form.getInputProps('Name')}
								/>

								{updateSettings.isError && (
									<Alert color="red" variant="light">
										Не удалось сохранить изменения
									</Alert>
								)}

								<Group justify="flex-end">
									<Button type="submit" loading={updateSettings.isPending}>
										Сохранить
									</Button>
								</Group>
							</Stack>
						</form>
					</Paper>
				</Stack>
			</Container>
		</>
	);
}
