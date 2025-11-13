import { useId, useRef } from "react";
import { Button } from "./ui/button";
import { Checkbox } from "./ui/checkbox";
import { Label } from "./ui/label";

function AuthForm() {
    const keep = useRef(false);
    const handleOAuth = (provider: string) => {
        window.location.href = `/api/users/auth/${provider}?keep=${keep.current}`;
    };

    const checkbox = useId();

    return (
        <>
            <Button onClick={() => handleOAuth("google")} variant="outline">
                Google
            </Button>
            <Button onClick={() => handleOAuth("github")} variant="outline">
                Github
            </Button>
            <div className="flex items-center gap-3">
                <Checkbox
                    id={checkbox}
                    onClick={() => {
                        keep.current = !keep.current;
                    }}
                />{" "}
                <Label htmlFor={checkbox}>Remember Me</Label>
            </div>
        </>
    );
}

export default AuthForm;
