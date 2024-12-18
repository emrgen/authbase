import React, {useState} from 'react';
import {Box, Text} from 'ink';
import Divider from 'ink-divider'
import {TextInput} from '@inkjs/ui';

type Props = {
	name: string | undefined;
};

export default function App({name = 'Stranger'}: Props) {
	const [greeting, setGreeting] = useState('Hello');
	return (
		<>
			<Divider title="Welcome to TUI" width={50}/>
			<Box>
				<TextInput onChange={setGreeting} placeholder="Type a greeting"/>
			</Box>

			<Box>
				<Text>
					{`${greeting}, ${name}!`}
				</Text>
			</Box>
		</>
	);
}
