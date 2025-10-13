import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { fetchCurrentUser, logout } from "@/services/api/users";

export const useCurrentUserQuery = () => {
    const queryClient = useQueryClient();

    const userQuery = useQuery({
        queryKey: ["me"],
        queryFn: fetchCurrentUser,
    });

    const logoutMutation = useMutation({
        mutationFn: logout,
        onSuccess: () => queryClient.setQueryData(["me"], null),
        retry: 2,
    });

    return {
        ...userQuery,
        logoutMutation,
    };
};
