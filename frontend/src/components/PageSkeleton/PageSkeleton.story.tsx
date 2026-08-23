import { StoryPreview } from '../StoryBook/StoryPreview';
import { PageSkeleton } from './PageSkeleton';

export default {
	title: 'PageSkeleton',
};

export function Default() {
	return (
		<StoryPreview isPage>
			<PageSkeleton />
		</StoryPreview>
	);
}
