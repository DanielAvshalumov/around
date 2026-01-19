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
    const [user, setIsUser] = useState({})

    const checkAuth = async () => {
        try {
            const url = 'http://localhost:8080/api/auth/me'
            const res = await axios.get(url,{withCredentials: true})
            if (res.status == 401) {
                console.log('Unauthorized')
                setIsAuthenticated(false)
                return
            }
            const userData = await res.data
            setIsAuthenticated(true)
            setIsUser(res.data)
            console.log(userData)
        } catch(err) {
            console.log("error", err)
        } finally {
            setIsLoading(false)
        }
    }

    useEffect(() => {
        checkAuth()
    },[])

    return (
        <AuthContext.Provider value={{ isAuthenticated, isLoading, user}}>
            {children}
        </AuthContext.Provider>
    )
}

export const useAuth = () => useContext(AuthContext);
