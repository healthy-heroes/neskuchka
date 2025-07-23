import { Header } from '@/components/Header/Header';
import { Features } from './Features/Features';
import { Intro } from './Intro/Intro';

export function LandingPage() {
	return (
		<>
			<Header />
			<Intro />
			<Features />
		</>
	);
}
