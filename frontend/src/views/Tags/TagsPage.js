import { Button } from "antd";
import { useNavigate } from "react-router-dom";



export default function TagsPage(params) {
    
    const navigateTo = useNavigate();

    return (
        <div className="tags">
            <Button
                type="primary"
                onClick={() => navigateTo("/dashboard/tags/add")}
            >创建标签</Button>  
            <h1>TagsPage</h1>
        </div>
    )
}