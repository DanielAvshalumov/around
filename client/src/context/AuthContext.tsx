'use client'

import axios from 'axios';
import { createContext, useContext, useEffect, useState} from 'react';
interface AuthContextType {
    isAuthenticated: boolean;
    isLoading: boolean;
    user: any | null;
}

const AuthContext = createContext<AuthContextType>({
    isAuthenticated: false,
    isLoading: true,
    user: null,
})

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [user, setUser] = useState({})

    const checkAuth = async () => {
        try {
            const url = 'http://localhost:8080/api/auth/me'
            const res = await axios.get(url,{withCredentials: true})
            
            const userData = await res.data
            setIsAuthenticated(true)
            setUser(userData)
            
        } catch(err: any) {
            console.log("error", err)
            if (err?.response.status == 401) {
                console.log('Unauthorized')
                console.log(isAuthenticated)
                setIsAuthenticated(false)
                
                if (window.location.href != "http://localhost:3000/") {
                    window.location.href = "/"
                }
                return
            }
            console.log(isAuthenticated)
        } finally {
            setIsLoading(false)
        }
    }

    useEffect(() => {
        checkAuth()
        console.log('after effect',isAuthenticated)
    },[])

    return (
        <AuthContext.Provider value={{ isAuthenticated, isLoading, user}}>
            {children}
        </AuthContext.Provider>
    )
}

export const useAuth = () => useContext(AuthContext);
