import { Button } from "./ui/button";

function AuthForm() {
    const handleOAuth = () => {
        window.location.href = "";
    };
    return (
        <Button onClick={() => handleOAuth()} variant="outline">
            Google
        </Button>
    );
}

export default AuthForm;
