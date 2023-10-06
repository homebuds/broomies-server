--- 356347e6-06e4-49e2-a31e-8c920851bbfd
INSERT INTO
        households (id, name)
VALUES
        ('356347e6-06e4-49e2-a31e-8c920851bbfd', 'Famjam');

--- 550e8400-e29b-41d4-a716-446655440000
--- f47ac10b-58cc-4372-a567-0e02b2c3d479
--- 6ba7b810-9dad-11d1-80b4-00c04fd430c8
--- 7b6452d0-c4cd-11e8-8bfc-4d7b37e1a3ab
--- a987fbc9-4bed-3078-cf07-9141ba07c9f3
INSERT INTO
        accounts (
                email,
                first_name,
                last_name,
                household_id,
                id,
                picture_url
        )
VALUES
        (
                'isaac.zhu@gmail.com',
                'Isaac',
                'Zhu',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                '550e8400-e29b-41d4-a716-446655440000',
                'https://www.nj.com/resizer/zovGSasCaR41h_yUGYHXbVTQW2A=/1280x0/smart/cloudfront-us-east-1.images.arcpublishing.com/advancelocal/SJGKVE5UNVESVCW7BBOHKQCZVE.jpg'
        ),
        (
                'gordon.wu@gmail.com',
                'Gordon',
                'Wu',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                'f47ac10b-58cc-4372-a567-0e02b2c3d479',
                'https://marketplace.canva.com/EAFewoMXU-4/1/0/1600w/canva-purple-pink-gradient-man-3d-avatar-0o0qE2T_kr8.jpg'
        ),
        (
                'rustam.nassyrov@gmail.com',
                'Rustam',
                'Nassyrov',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
                'https://w7.pngwing.com/pngs/114/914/png-transparent-rust-programming-language-logo-machine-learning-haskell-crab-animals-cartoon-crab.png'
        ),
        (
                'dennis.li@gmail.com',
                'Dennis',
                'Li',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                '7b6452d0-c4cd-11e8-8bfc-4d7b37e1a3ab',
                'https://i.pinimg.com/474x/1e/29/cc/1e29cc76516f3f38cc8589f4ce831ee1.jpg'
        );

INSERT INTO
        roommate_scores (id, account_id, household_id, points)
VALUES
        (
                uuid_generate_v4(),
                '550e8400-e29b-41d4-a716-446655440000',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                500
        ),
        (
                uuid_generate_v4(),
                'f47ac10b-58cc-4372-a567-0e02b2c3d479',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                550
        ),
        (
                uuid_generate_v4(),
                '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                650
        ),
        (
                uuid_generate_v4(),
                '7b6452d0-c4cd-11e8-8bfc-4d7b37e1a3ab',
                '356347e6-06e4-49e2-a31e-8c920851bbfd',
                900
        );

-- {
--     accountId: "a1",
--     firstName: "John",
--     lastName: "Doe",
--     email: "john.doe@example.com",
--     photo: "https://www.nj.com/resizer/zovGSasCaR41h_yUGYHXbVTQW2A=/1280x0/smart/cloudfront-us-east-1.images.arcpublishing.com/advancelocal/SJGKVE5UNVESVCW7BBOHKQCZVE.jpg"
-- },
-- {
--     accountId: "a2",
--     firstName: "Jane",
--     lastName: "Doe",
--     email: "jane.doe@example.com",
--     photo: "https://people.com/thmb/84W5-9FnCb0XLaqaoYwHasY5GwI=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc():focal(216x0:218x2)/robert-pattinson-435-2-3f3472a03106439abee37574a6b8cef7.jpg"
-- },
-- (
--         'John',
--         'Doe',
--         'john.doe@gmail.com',
--         'https://www.nj.com/resizer/zovGSasCaR41h_yUGYHXbVTQW2A=/1280x0/smart/cloudfront-us-east-1.images.arcpublishing.com/advancelocal/SJGKVE5UNVESVCW7BBOHKQCZVE.jpg'
-- ) (
--         'Jane',
--         'Doe',
--         'jane.doe@gmail.com',
--         'https://people.com/thmb/84W5-9FnCb0XLaqaoYwHasY5GwI=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc():focal(216x0:218x2)/robert-pattinson-435-2-3f3472a03106439abee37574a6b8cef7.jpg'
-- )