import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { LandingPage } from './Landing.page';

export default {
	title: 'Pages/Landing',
};

export function Default() {
	return (
		<StoryPreview isPage>
			<LandingPage />
		</StoryPreview>
	);
}
