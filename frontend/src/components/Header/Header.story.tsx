import { createApiServiceMock } from '@/api/fixtures/api';
import { createAuthServiceMock, mockUser } from '@/api/fixtures/auth';
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
		auth: createAuthServiceMock({ user: mockUser }),
	});

	return (
		<StoryPreview apiService={apiService}>
			<Header />
		</StoryPreview>
	);
}
