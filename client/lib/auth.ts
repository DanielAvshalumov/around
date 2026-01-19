import axios from "axios";

export interface User {
    id: string;
    email: string;
    name: string;
    provider: string;
}

export class AuthService {

    static async login() {
        console.log("logging in")
        console.log(process.env.NEXT_PUBLIC_GOOGLE_OAUTH_CLIENT_ID)
        const GoogleOAuthUrl = 'https://accounts.google.com/o/oauth2/v2/auth'
        const payload = {
            client_id :process.env.NEXT_PUBLIC_GOOGLE_OAUTH_CLIENT_ID as string,
            redirect_uri :process.env.NEXT_PUBLIC_OAUTH_REDIRECT_URI as string,
            scope: "email",
            response_type: "code",
            state: process.env.NEXT_PUBLIC_GOOGLE_OAUTH_SECRET_STATE as string,
        }
      
        const params = new URLSearchParams(payload).toString()
        const res = `${GoogleOAuthUrl}?${params}`
        
        window.location.href = res
    }

    static async logout() {
        
    }


    static async verifyUser(code: string, provider: string) {

        window.location.href = `${process.env.NEXT_PUBLIC_OAUTH_REDIRECT_URI}?code=${code}&provider=${provider}`

    }
}