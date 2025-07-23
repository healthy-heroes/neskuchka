import { Button, Container, Text, Title } from '@mantine/core';
import classes from './Intro.module.css';

/**
 * 
 * Сделайте спорт частью своей жизни с Нескучным Спортом
 * Простые упражнения, которые займут 10 минут и помогут взбодриться.
Оборудование не нужно.
 */

export function Intro() {
	return (
		<div className={classes.root}>
			<Container size="lg">
				<div className={classes.inner}>
					<div className={classes.content}>
						<Title className={classes.title}>Сделайте спорт частью своей жизни</Title>
						<Text className={classes.description} mt={30}>
							Простые упражнения, которые займут 10 минут и помогут взбодриться. Оборудование не
							нужно.
						</Text>
						<Button variant="gradient" size="xl" className={classes.control} mt={40}>
							Присоединиться
						</Button>
					</div>
				</div>
			</Container>
		</div>
	);
}
