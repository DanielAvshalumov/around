import axios from 'axios'

export default async function getBacklinkDomains(industry: string) {
    try {
        const res = await axios.get(`http://localhost:8000/industry_domains?industry=${industry}`)
        const data = await res.data
        const response = data.split('\n')
        return response
    } catch(error: any) {
        console.log("GET\tHTTP/1.1\t/industry_domains", error)
    }
}