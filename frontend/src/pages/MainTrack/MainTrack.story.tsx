import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { MainTrackPage } from './MainTrack.page';

export default {
	title: 'Pages/MainTrack',
};

export function Default() {
	return (
		<StoryPreview isPage>
			<MainTrackPage />
		</StoryPreview>
	);
}
