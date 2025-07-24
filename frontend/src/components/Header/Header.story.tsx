import { StoryPreview } from '../StoryBook/StoryPreview';
import { Header } from './Header';

export default {
	title: 'Header',
	viewport: {
		options: {
			kindleFire2: {
				name: 'Kindle Fire 2',
				styles: { width: '600px', height: '963px' },
			},
		},
	},
};

export function Default() {
	return (
		<StoryPreview>
			<Header />
		</StoryPreview>
	);
}

// Дополнительная история специально для мобильного вида
export function Mobile() {
	return (
		<StoryPreview>
			<Header />
		</StoryPreview>
	);
}

Mobile.parameters = {
	viewport: {
		options: {
			kindleFire2: {
				name: 'Kindle Fire 2',
				styles: { width: '600px', height: '963px' },
			},
		},
	},
};
