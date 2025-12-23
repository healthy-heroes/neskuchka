import { useContext } from 'react';
import { ApiContext } from './provider';

export function useApi() {
	const context = useContext(ApiContext);
	if (context === null) {
		throw new Error('useApi must be used within an ApiProvider');
	}

	return context;
}
