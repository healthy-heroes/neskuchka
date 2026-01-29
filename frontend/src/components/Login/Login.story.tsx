import type { Meta, StoryObj } from '@storybook/react';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { LoginForm } from './LoginForm';
import { LoginSuccess } from './LoginSuccess';

const meta: Meta<typeof LoginForm> = {
	title: 'Auth/Login',
	component: LoginForm,
	parameters: {
		layout: 'centered',
	},
	tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof LoginForm>;

// Default state - empty form
export const Default: Story = {
	render: () => (
		<StoryPreview>
			<LoginForm
				onSubmit={(email: string) => {
					console.log('Login submitted with email:', email);
				}}
				isSubmitting={false}
				error={null}
			/>
		</StoryPreview>
	),
};

// Loading state - form is submitting
export const Loading: Story = {
	render: () => (
		<StoryPreview>
			<LoginForm
				onSubmit={(email: string) => {
					console.log('Login submitted with email:', email);
				}}
				isSubmitting
				error={null}
			/>
		</StoryPreview>
	),
};

// Error state - server error
export const WithError: Story = {
	render: () => (
		<StoryPreview>
			<LoginForm
				onSubmit={(email: string) => {
					console.log('Login submitted with email:', email);
				}}
				error={new Error('Не удалось отправить письмо. Попробуйте позже.')}
			/>
		</StoryPreview>
	),
};

// Success state - email sent
export const Success: Story = {
	render: () => (
		<StoryPreview>
			<LoginSuccess />
		</StoryPreview>
	),
};
