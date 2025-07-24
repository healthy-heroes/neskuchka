import { Paper, PaperProps } from '@mantine/core';

export interface StoryPreviewProps {
	children: React.ReactNode;

	paperOptions?: PaperProps;
}

export function StoryPreview({ children, paperOptions }: StoryPreviewProps) {
	return (
		<Paper shadow="xs" p="sm" m="sm" {...paperOptions}>
			{children}
		</Paper>
	);
}
