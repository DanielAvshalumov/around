import axios from 'axios'

type Payload = {
    industry: string;
    comp_domains: string[] | null;
    browser: string;
}

export default async function getForumProductLinks(payload?: Payload) {
    try {
        // const res = await axios.post("http://localhost:8080/back-link",{"comp_domains":["amazon.com","ajmadison.com","homedepot.com","bestbuy.com","build.com","lowes.com"],"industry":`${payload?.industry}`,"browser":"duckduckgo"},{'headers':{'Content-Type' : 'application/json'}})
        const res = await axios.post("http://localhost:8080/back-link",{"industry":`${payload?.industry}`,"browser":"duckduckgo"},{'headers':{'Content-Type' : 'application/json'}})
        const data = await res.data
        return data;
    } catch (error: any) {
        console.log("error",error)
    }
}