import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { ButtonSpecimen } from './ButtonSpecimen';

export default {
	title: 'Primitives/Buttons',
};

export function Default() {
	return (
		<StoryPreview>
			<ButtonSpecimen />
		</StoryPreview>
	);
}
