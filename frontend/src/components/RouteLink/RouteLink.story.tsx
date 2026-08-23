import { Group } from '@mantine/core';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { RouteLink } from './RouteLink';

export default {
	title: 'RouteLink',
};

export function Default() {
	return (
		<StoryPreview>
			<Group gap="lg">
				<RouteLink to="/">Обычная ссылка</RouteLink>
				<RouteLink to="/" underline="never">
					Без подчёркивания
				</RouteLink>
				<RouteLink to="/" c="gray.7">
					Вторичная
				</RouteLink>
			</Group>
		</StoryPreview>
	);
}
