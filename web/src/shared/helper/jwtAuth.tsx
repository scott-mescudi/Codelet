import {jwtDecode} from 'jwt-decode'

export function IsTokenExpired(token: string): boolean {
	try {
		const decoded: {exp?: number} = jwtDecode(token)
		const now = Date.now() / 1000

		return decoded.exp !== undefined ? decoded.exp < now : true
	} catch (error) {
		return true
	}
}
