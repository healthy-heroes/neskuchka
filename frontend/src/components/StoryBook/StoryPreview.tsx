import { Paper, PaperProps } from '@mantine/core';

export interface StoryPreviewProps {
	children: React.ReactNode;

	paperOptions?: PaperProps;
	isPage?: boolean;
}

export function StoryPreview({ children, paperOptions, isPage }: StoryPreviewProps) {
	const defaultOptions: PaperProps = {
		shadow: 'xs',
		m: 'sm',
	};

	if (isPage) {
		defaultOptions.bd = '1px solid var(--mantine-color-gray-2)';
	} else {
		defaultOptions.p = 'sm';
	}

	return (
		<Paper {...defaultOptions} {...paperOptions}>
			{children}
		</Paper>
	);
}
