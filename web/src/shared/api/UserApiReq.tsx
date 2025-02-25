interface LoginRequest {
	email: string
	password: string
}

interface LoginResponse {
	access_token: string
}

interface SignupRequest {
	username: string
	email: string
	password: string
	role: string
}

interface ErrorResp {
	error: string
	code: number
}

interface GetUsernameResponse {
	username: string
}

export async function Signup(
	username: string,
	email: string,
	password: string
): Promise<number> {
	const data: SignupRequest = {
		username,
		email,
		password,
		role: 'user'
	}

	try {
		const resp = await fetch('http://localhost:3021/api/v1/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		})

		if (!resp.ok) {
			return resp.status
		}

		return 200
	} catch (err) {
		return 500
	}
}

export async function Login(email: string, password: string): Promise<number> {
	const data: LoginRequest = {
		email,
		password
	}

	try {
		const resp = await fetch('http://localhost:3021/api/v1/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		})

		if (!resp.ok) {
			return resp.status
		}

		const token = (await resp.json()) as LoginResponse
		localStorage.setItem('ACCESS_TOKEN', token.access_token)

		return 200
	} catch (err) {
		return 500
	}
}

export async function Logout(token: string) {
	if (token === '') return

	try {
		const resp = await fetch(`http://localhost:3021/api/v1/logout`, {
			method: 'POST',
			headers: {
				Authorization: token
			}
		})

		if (!resp.ok) {
			const errResp = (await resp.json()) as ErrorResp
			console.log(errResp)
			return
		}
	} catch (err) {
		console.error(err)
		return
	}
}


export async function GetUsername(token: string): Promise<string | undefined> {
	if (token === '') return undefined

	try {
		const resp = await fetch(`http://localhost:3021/api/v1/username`, {
			method: 'GET',
			headers: {
				Authorization: token
			}
		})

		if (!resp.ok) {
			const errResp = (await resp.json()) as ErrorResp
			console.log(errResp)
			return undefined
		}
		
		const info = (await resp.json()) as GetUsernameResponse

		return info.username
	} catch (err) {
		console.error(err)
		return undefined
	}
}
