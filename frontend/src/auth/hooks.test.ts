import { act, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createApiServiceMock } from '@/api/fixtures/api';
import { createAuthServiceMock, mockUser } from '@/api/fixtures/auth';
import { useApi } from '@/api/hooks';
import { AuthKeys } from '@/api/services/auth';
import { createTestQueryClient, renderHook } from '../../test-utils';
import { useAuth, useIsOwner } from './hooks';

// Mock useApi hook
vi.mock('@/api/hooks', () => ({
	useApi: vi.fn(),
}));

describe('useAuth', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('should return loading state initially', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useAuth());

		expect(result.current.isLoading).toBe(true);
		expect(result.current.user).toBe(null);
		expect(result.current.isAuthenticated).toBe(false);
	});

	it('should return user when authenticated', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useAuth());

		await waitFor(() => {
			expect(result.current.isLoading).toBe(false);
		});

		expect(result.current.user).toEqual(mockUser);
		expect(result.current.isAuthenticated).toBe(true);
	});

	it('should return null user when not authenticated', async () => {
		const authMock = createAuthServiceMock({ user: null });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useAuth());

		await waitFor(() => {
			expect(result.current.isLoading).toBe(false);
		});

		expect(result.current.user).toBe(null);
		expect(result.current.isAuthenticated).toBe(false);
	});

	it('should clear user data on logout', async () => {
		const queryClient = createTestQueryClient();

		// Pre-fill cache with user data
		queryClient.setQueryData(AuthKeys.user, { data: mockUser });

		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useAuth(), { queryClient });

		// Wait for initial data to be available
		await waitFor(() => {
			expect(result.current.user).toEqual(mockUser);
		});
		expect(result.current.isAuthenticated).toBe(true);

		// Perform logout
		await act(async () => {
			await result.current.logout();
		});

		// User should be cleared
		await waitFor(() => {
			expect(result.current.user).toBe(null);
		});
		expect(result.current.isAuthenticated).toBe(false);
	});

	it('should handle logout error gracefully', async () => {
		const queryClient = createTestQueryClient();
		queryClient.setQueryData(AuthKeys.user, { data: mockUser });

		const logoutError = new Error('Network error');
		const authMock = createAuthServiceMock({ user: mockUser, logoutError });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useAuth(), { queryClient });

		// Logout should throw
		await expect(
			act(async () => {
				await result.current.logout();
			})
		).rejects.toThrow('Network error');

		// User should still be there (logout failed)
		expect(result.current.user).toEqual(mockUser);
		expect(result.current.isAuthenticated).toBe(true);
	});
});

describe('useIsOwner', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('should return true when user is owner', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useIsOwner(mockUser.ID));

		await waitFor(() => {
			expect(result.current).toBe(true);
		});
	});

	it('should return false when user is not owner', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useIsOwner('different-user-id'));

		await waitFor(() => {
			expect(result.current).toBe(false);
		});
	});

	it('should return false when user is not authenticated', async () => {
		const authMock = createAuthServiceMock({ user: null });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useIsOwner('some-owner-id'));

		await waitFor(() => {
			expect(result.current).toBe(false);
		});
	});

	it('should return false when ownerID is undefined', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useIsOwner(undefined));

		expect(result.current).toBe(false);
	});

	it('should return false when ownerID is empty string', async () => {
		const authMock = createAuthServiceMock({ user: mockUser });
		vi.mocked(useApi).mockReturnValue(createApiServiceMock({ auth: authMock }));

		const { result } = renderHook(() => useIsOwner(''));

		expect(result.current).toBe(false);
	});
});
