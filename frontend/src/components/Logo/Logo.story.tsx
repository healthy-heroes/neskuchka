import { StoryPreview } from '../StoryBook/StoryPreview';
import { Logo } from './Logo';

export default {
	title: 'Logo',
};

export function Default() {
	return (
		<StoryPreview>
			<Logo />
		</StoryPreview>
	);
}
