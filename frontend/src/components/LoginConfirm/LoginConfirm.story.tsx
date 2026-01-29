import type { Meta, StoryObj } from '@storybook/react';
import { createApiServiceMock } from '@/api/fixtures/api';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { LoginConfirm } from './LoginConfirm';

const meta: Meta<typeof LoginConfirm> = {
	title: 'Auth/LoginConfirm',
	component: LoginConfirm,
	parameters: {
		layout: 'centered',
	},
	tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof LoginConfirm>;

// Loading state - confirming token
export const Loading: Story = {
	render: () => {
		const apiService = createApiServiceMock({
			auth: {
				confirmLoginMutation: () => ({
					mutationFn: () => new Promise(() => {}), // Never resolves to stay in loading
				}),
			},
		});

		return (
			<StoryPreview apiService={apiService}>
				<LoginConfirm token="valid-token-123" />
			</StoryPreview>
		);
	},
};

// Success state - login confirmed
export const Success: Story = {
	render: () => {
		const apiService = createApiServiceMock({
			auth: {
				confirmLoginMutation: () => ({
					mutationFn: () => Promise.resolve({ success: true }),
				}),
			},
		});

		return (
			<StoryPreview apiService={apiService}>
				<LoginConfirm token="valid-token-123" />
			</StoryPreview>
		);
	},
};

// No token state
export const NoToken: Story = {
	render: () => (
		<StoryPreview>
			<LoginConfirm />
		</StoryPreview>
	),
};

// Error state - invalid token
export const WithError: Story = {
	render: () => {
		const apiService = createApiServiceMock({
			auth: {
				confirmLoginMutation: () => ({
					mutationFn: () => Promise.reject(new Error('Токен недействителен или истёк')),
				}),
			},
		});

		return (
			<StoryPreview apiService={apiService}>
				<LoginConfirm token="invalid-token" />
			</StoryPreview>
		);
	},
};
