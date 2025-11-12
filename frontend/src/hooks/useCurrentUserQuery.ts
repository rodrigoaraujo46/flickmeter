import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { fetchCurrentUser, logout } from "@/services/api/users";

export const useCurrentUserQuery = () => {
    const queryClient = useQueryClient();

    const userQuery = useQuery({
        queryKey: ["users", "me"],
        queryFn: fetchCurrentUser,
    });

    const logoutMutation = useMutation({
        mutationFn: logout,
        onSuccess: (ok) => {
            if (ok) {
                queryClient.resetQueries();
            }
        },
    });

    return {
        ...userQuery,
        logoutMutation,
    };
};
