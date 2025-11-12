import * as React from "react";

interface ErrorBoundaryProps {
    fallback: React.ReactNode;
    children: React.ReactNode;
}

interface ErrorBoundaryState {
    hasError: boolean;
    error?: Error;
}

export class ErrorBoundary extends React.Component<
    ErrorBoundaryProps,
    ErrorBoundaryState
> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = { hasError: false };
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return { hasError: true, error: error };
    }

    componentDidCatch(error: Error, info: React.ErrorInfo) {
        // Replace this with your actual error logging function
        console.error("Caught error:", error, info.componentStack);
    }

    render() {
        if (this.state.hasError) {
            return (
                <>
                    <div>{this.props.fallback}</div>
                    <div>{this.state.error?.message}</div>{" "}
                </>
            );
        }

        return this.props.children;
    }
}
