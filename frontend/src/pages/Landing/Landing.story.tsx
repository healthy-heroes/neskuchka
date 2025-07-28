import { PaperProps } from '@mantine/core';
import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { LandingPage } from './Landing.page';

export default {
	title: 'Pages/Landing',
};

export function Default() {
	const paperOptions: PaperProps = {
		bd: '1px solid var(--mantine-color-gray-2)',
		p: '0',
	};

	return (
		<StoryPreview paperOptions={paperOptions}>
			<LandingPage />
		</StoryPreview>
	);
}
