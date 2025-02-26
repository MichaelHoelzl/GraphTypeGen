import { gql } from '@apollo/client';

export interface User {
	id: string;
	name: string;
	email: string;
	password: string;
	createdAt: string;
	updatedAt: string;
}

async function createUser(name: string, email: string, password: string): Promise<[User | null, any | null]> {
	try {
		const response = await apolloClient.mutate({
			mutation: gql(`
				 mutation createUser($name: String!, $email: String!, $password: String!) {
					createUser(name: $name, email: $email, password: $password) {
						id 
						name 
						email 
						password 
						createdAt 
						updatedAt 
					}
				}
				`),
			variables: {
								name: name,
				email: email,
				password: password
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.createUser) {
			return [null, new Error(`Invalid response structure`)];
		}
	
		return [response.data.createUser, null];
	} catch (error) {
		return [null, error];
	}
}

async function deleteUser(id: string): Promise<[User | null, any | null]> {
	try {
		const response = await apolloClient.mutate({
			mutation: gql(`
				 mutation deleteUser($id: ID!) {
					deleteUser(id: $id) {
						id 
						name 
						email 
						password 
						createdAt 
						updatedAt 
					}
				}
				`),
			variables: {
								id: id
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.deleteUser) {
			return [null, new Error(`Invalid response structure`)];
		}
	
		return [response.data.deleteUser, null];
	} catch (error) {
		return [null, error];
	}
}

async function updateUser(id: string, name: string, email: string, password: string): Promise<[User | null, any | null]> {
	try {
		const response = await apolloClient.mutate({
			mutation: gql(`
				 mutation updateUser($id: ID!, $name: String, $email: String, $password: String) {
					updateUser(id: $id, name: $name, email: $email, password: $password) {
						id 
						name 
						email 
						password 
						createdAt 
						updatedAt 
					}
				}
				`),
			variables: {
								id: id,
				name: name,
				email: email,
				password: password
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.updateUser) {
			return [null, new Error(`Invalid response structure`)];
		}
	
		return [response.data.updateUser, null];
	} catch (error) {
		return [null, error];
	}
}

async function user(id: string): Promise<[User | null, any | null]> {
	try {
		const response = await apolloClient.query({
			query: gql(`
				 query user($id: ID!) {
					user(id: $id) {
						id 
						name 
						email 
						password 
						createdAt 
						updatedAt 
					}
				}
				`),
			variables: {
								id: id
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.user) {
			return [null, new Error(`Invalid response structure`)];
		}
	
		return [response.data.user, null];
	} catch (error) {
		return [null, error];
	}
}

async function users(): Promise<[User[] | null, any | null]> {
	try {
		const response = await apolloClient.query({
			query: gql(`
				 query users {
					users {
						id 
						name 
						email 
						password 
						createdAt 
						updatedAt 
					}
				}
				`),
			variables: {
				
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.users) {
			return [null, new Error(`Invalid response structure`)];
		}
	
		return [response.data.users, null];
	} catch (error) {
		return [null, error];
	}
}

export const query = {
	user,
	users
};
export const mutation = {
	createUser,
	deleteUser,
	updateUser
};


