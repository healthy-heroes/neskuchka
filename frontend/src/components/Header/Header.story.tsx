import { createApiServiceMock } from '@/api/fixtures/api';
import avatar from '@/api/fixtures/avatar.jpg';
import { createUserServiceMock, mockUser } from '@/api/fixtures/user';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Header } from './Header';

export default {
	title: 'Header',
};

export function LoggedOut() {
	return (
		<StoryPreview>
			<Header />
		</StoryPreview>
	);
}

export function LoggedIn() {
	const apiService = createApiServiceMock({
		user: createUserServiceMock({ user: mockUser }),
	});

	return (
		<StoryPreview apiService={apiService}>
			<Header />
		</StoryPreview>
	);
}

export function LoggedInWithAvatar() {
	const apiService = createApiServiceMock({
		user: createUserServiceMock({
			user: { ...mockUser, Avatar: avatar },
		}),
	});

	return (
		<StoryPreview apiService={apiService}>
			<Header />
		</StoryPreview>
	);
}
