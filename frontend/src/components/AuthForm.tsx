import { useRef } from "react";
import { Button } from "./Button";

function AuthForm() {
    const keep = useRef(false);
    const handleOAuth = (provider: string) => {
        window.location.href = `/api/user/auth/${provider}?keep=${keep.current}`;
    };

    return (
        <>
            <Button onClick={() => handleOAuth("google")} variant="outline">
                Google
            </Button>
            <Button onClick={() => handleOAuth("github")} variant="outline">
                Github
            </Button>

            <label>
                <input
                    name="keep-me"
                    type="checkbox"
                    onClick={() => {
                        keep.current = !keep.current;
                    }}
                />{" "}
                Remember Me
            </label>
        </>
    );
}

export default AuthForm;
