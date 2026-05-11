"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useSetAtom } from "jotai";
import { Box, Button, Heading, Input, Text, VStack } from "@chakra-ui/react";
import { useLogin } from "@/api/generated/ninjo";
import { userAtom } from "@/stores/auth";
import Link from "next/link";

export default function LoginPage() {
  const router = useRouter();
  const setUser = useSetAtom(userAtom);
  const loginMutation = useLogin();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    loginMutation.mutate(
      { data: { email, password } },
      {
        onSuccess: (user) => {
          setUser(user);
          router.push("/");
        },
        onError: (err) => {
          setError(err instanceof Error ? err.message : "ログインに失敗しました");
        },
      }
    );
  };

  return (
    <Box as="form" onSubmit={handleSubmit} w="full" bg="white" p={8} borderRadius="lg" shadow="sm">
      <VStack gap={4}>
        <Heading size="lg">ログイン</Heading>

        {error && (
          <Text color="red.500" fontSize="sm">
            {error}
          </Text>
        )}

        <VStack gap={3} w="full">
          <Input
            type="email"
            placeholder="メールアドレス"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <Input
            type="password"
            placeholder="パスワード"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </VStack>

        <Button
          type="submit"
          w="full"
          colorPalette="blue"
          loading={loginMutation.isPending}
        >
          ログイン
        </Button>

        <Text fontSize="sm" color="gray.600">
          アカウントをお持ちでない方は{" "}
          <Link href="/signup" style={{ color: "var(--chakra-colors-blue-500)" }}>
            新規登録
          </Link>
        </Text>
      </VStack>
    </Box>
  );
}
